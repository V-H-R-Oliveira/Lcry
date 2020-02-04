# Lcry - A Golang toy malware

- The Lcry malware is a toy malware that installs two services in your systemd.
- The services grant persistence after reboot, and they will init with the system startup.
- The lcry.service is the main service, which releases a ransomware component to encrypt the files in the $HOME.
- The monitor.service is a service designed to launch the Lcry process monitor.

## How the attack works:
- After the execution of the loader, it will create a hidden dir in your /, register both of the services, start a tor instance, and launch the ransomware component.
- The ransomware component will get an encryption key from the server and it will encrypt all the files with AES256-CTR, authenticated with HMAC-SHA256.
- After the encryption process, the ransomware executes the built-in monitor to monitoring the monitor.lcry.  
- The ransomware will enter in a stand-by mode.
- The lcry.service monitor and the monitor.service monitor will send the current status of each process in a specific time interval to the command server.
- All the communication between the command server and the malware will use a Tor hidden service.