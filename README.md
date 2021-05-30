#go_NetworkMonitoring

This was created for my home lab, monitors my services and VM's

<h1>Features:</h1>
<ul>
    <li>Scans IP's given from a hosts.csv file</li>
    <li>Multithreading ensures all IP's get scanned equally</li>
    <li>Will report if a host does not respond to an ICMP packet for 10 seconds</li>
    <li>Will report if a device has more than 35% packet loss</li>
    <li>Reports come in the form of an email alert.</li>
</ul>