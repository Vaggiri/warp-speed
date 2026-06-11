# Privacy Policy for Warp-Speed

**Effective Date:** June 11, 2026

## Overview
Warp-Speed is an open-source Command Line Interface (CLI) tool for network diagnostics. We are committed to protecting your privacy. This policy outlines how Warp-Speed handles data.

## Data Collection
Warp-Speed **does not** collect, store, or transmit any personally identifiable information (PII), analytics, or usage telemetry to its developers or any third parties.

## Network Activity
To function correctly, the application performs the following network activities directly from your local machine:
- Contacts public DNS servers (e.g., Cloudflare 1.1.1.1, Google 8.8.8.8) to measure latency.
- Connects to global Speedtest nodes to measure download and upload bandwidth.
- Fetches a list of available speed test servers via public APIs.

All speed test results and historical data are stored **locally** on your device (typically in `~/.warp-speed/history.json`) so you can view your past performance. This data is never uploaded to any remote server.

## Changes to this Policy
We may update this Privacy Policy from time to time as features are added. Any changes will be reflected directly in our GitHub repository.

## Contact
If you have any questions or concerns about this privacy policy, please open an issue on our GitHub repository.
