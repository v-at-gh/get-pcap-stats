package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	dirPath     string
	statsTypes  string
	suffix      string
	overwrite   bool
	yes         bool
	parallelism int
)

func init() {
	flag.StringVar(&dirPath, "dir", ".", "Directory to search for .pcap or .pcapng files")
	flag.StringVar(&statsTypes, "stats", "", "Statistics types (space-separated) or path to file containing types")
	flag.StringVar(&suffix, "suffix", ".total-stats.txt", "Suffix for resulting statistics files")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite existing statistics files")
	flag.BoolVar(&yes, "yes", false, "Process files without asking for confirmation")
	flag.IntVar(&parallelism, "workers", runtime.NumCPU(), "Number of files to process in parallel")
}

func main() {
	flag.Parse()

	files, err := findPcapFiles(dirPath, overwrite, suffix)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No pcap or pcapng files need processing.")
		return
	}

	displayFiles(files)

	if !yes {
		if !confirm("Do you want to proceed? (y/n): ") {
			fmt.Println("Aborted.")
			return
		}
	}

	args, err := buildArgs(statsTypes)
	if err != nil {
		fmt.Printf("Error building arguments: %v\n", err)
		os.Exit(1)
	}

	processFiles(files, args)
}

func findPcapFiles(dir string, overwrite bool, suffix string) ([]string, error) {
	var filesToProcess []string
	err := filepath.WalkDir(
		dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (strings.HasSuffix(path, ".pcap") || strings.HasSuffix(path, ".pcapng")) {
				statsPath := strings.Replace(path, filepath.Ext(path), suffix, 1)
				if overwrite || !fileExists(statsPath) {
					filesToProcess = append(filesToProcess, path)
				}
			}
			return nil
		},
	)
	return filesToProcess, err
}

func displayFiles(files []string) {
	fmt.Println("The following files will be processed:")
	for i, f := range files {
		fmt.Printf("%d: %s\n", i+1, f)
	}
	fmt.Printf("Total %d files.\n", len(files))
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		return text == "y" || text == "Y"
	}
	return false
}

func processFiles(files []string, args []string) {
	sem := make(chan struct{}, parallelism)
	var wg sync.WaitGroup
	totalFiles := len(files)

	for i, file := range files {
		wg.Add(1)
		sem <- struct{}{}

		go func(f string, idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			processFile(f, idx, totalFiles, args)
		}(file, i)
	}

	wg.Wait()
	close(sem)
}

func processFile(filePath string, index int, totalFiles int, args []string) {
	statsPath := strings.Replace(filePath, filepath.Ext(filePath), suffix, 1)
	fmt.Printf("Processing %d/%d: %s...\n", index+1, totalFiles, filePath)

	if err := extractStats(filePath, statsPath, args); err != nil {
		fmt.Printf("Error processing %s: %v\n", filePath, err)
	}
}

func buildArgs(statsOption string) ([]string, error) {
	args := []string{"-q"}
	var stats []string

	if statsOption == "" {
		stats = []string{
			"afp,srt", "ancp,tree",
			"ansi_a,bsmap", "ansi_a,dtap", "ansi_map",
			"asap,stat", "bacapp_instanceid,tree",
			"bacapp_ip,tree", "bacapp_objectid,tree",
			"bacapp_service,tree", "calcappprotocol,stat",
			"camel,counter", "camel,srt", "collectd,tree",
			"componentstatusprotocol,stat", "conv,bluetooth",
			"conv,bpv7", "conv,dccp", "conv,eth",
			"conv,fc", "conv,fddi", "conv,ip",
			"conv,ipv6", "conv,ipx", "conv,jxta",
			"conv,ltp", "conv,mptcp", "conv,ncp",
			"conv,opensafety", "conv,rsvp", "conv,sctp",
			"conv,sll", "conv,tcp", "conv,tr",
			"conv,udp", "conv,usb", "conv,wlan", "conv,wpan",
			"conv,zbee_nwk", "credentials", "dests,tree",
			"dhcp,stat", "diameter,avp", "diameter,srt",
			"dns,tree", "dns_qr,tree", "e2ap,tree",
			"endpoints,bluetooth", "endpoints,bpv7",
			"endpoints,dccp", "endpoints,eth", "endpoints,fc",
			"endpoints,fddi", "endpoints,ip", "endpoints,ipv6",
			"endpoints,ipx", "endpoints,jxta", "endpoints,ltp",
			"endpoints,mptcp", "endpoints,ncp", "endpoints,opensafety",
			"endpoints,rsvp", "endpoints,sctp", "endpoints,sll",
			"endpoints,tcp", "endpoints,tr", "endpoints,udp",
			"endpoints,usb", "endpoints,wlan", "endpoints,wpan",
			"endpoints,zbee_nwk", "enrp,stat", "expert", "f1ap,tree",
			"f5_tmm_dist,tree", "f5_virt_dist,tree", "fc,srt",
			"fractalgeneratorprotocol,stat", "gsm_a", "gsm_a,bssmap",
			"gsm_a,dtap_cc", "gsm_a,dtap_gmm", "gsm_a,dtap_mm",
			"gsm_a,dtap_rr", "gsm_a,dtap_sacch", "gsm_a,dtap_sm",
			"gsm_a,dtap_sms", "gsm_a,dtap_ss", "gsm_a,dtap_tp",
			"gsm_map,operation", "gtp,srt", "gtpv2,srt",
			"h225,counter", "h225_ras,rtd", "hart_ip,tree",
			"hosts", "hpfeeds,tree", "http,stat", "http,tree",
			"http2,tree", "http_req,tree", "http_seq,tree",
			"http_srv,tree", "icmp,srt", "icmpv6,srt",
			"io,phs", "ip_hosts,tree", "ip_srcdst,tree",
			"ip_ttl,tree", "ipv6_dests,tree", "ipv6_hop,tree",
			"ipv6_hosts,tree", "ipv6_ptype,tree", "ipv6_srcdst,tree",
			"isup_msg,tree", "kerberos,srt", "lbmr_queue_ads_queue,tree",
			"lbmr_queue_ads_source,tree", "lbmr_queue_queries_queue,tree",
			"lbmr_queue_queries_receiver,tree", "lbmr_topic_ads_source,tree",
			"lbmr_topic_ads_topic,tree", "lbmr_topic_ads_transport,tree",
			"lbmr_topic_queries_pattern,tree", "lbmr_topic_queries_pattern_receiver,tree",
			"lbmr_topic_queries_receiver,tree", "lbmr_topic_queries_topic,tree",
			"ldap,srt", "ltp,tree", "mac-3gpp,stat", "mgcp,rtd",
			"mtp3,msus", "ncp,srt", "nfsv4,srt", "ngap,tree",
			"npm,stat", "osmux,tree", "pfcp,srt", "pingpongprotocol,stat",
			"plen,tree", "ptype,tree", "radius,rtd", "rlc-3gpp,stat",
			"rtp,streams", "rtsp,stat", "rtsp,tree", "sametime,tree",
			"sctp,stat", "sip,stat", "smpp_commands,tree", "snmp,srt",
			"someip_messages,tree", "someipsd_entries,tree", "ssprotocol,stat",
			"sv", "ucp_messages,tree", "wsp,stat",
		}
	} else if fileExists(statsOption) {
		file, err := os.Open(statsOption)
		if err != nil {
			return nil, fmt.Errorf("failed to open stats file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			stats = append(stats, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read stats file: %w", err)
		}
	} else {
		stats = strings.Fields(statsOption)
	}

	for _, stat := range stats {
		args = append(args, "-z", stat)
	}

	return args, nil
}

func extractStats(pcapPath, statsPath string, args []string) error {
	cmdArgs := append([]string{"-r", pcapPath}, args...)

	cmd := exec.Command("tshark", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	if err := os.WriteFile(statsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output to stats file: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
