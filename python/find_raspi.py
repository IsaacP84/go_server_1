from scapy.all import srp, Ether, ARP
import sys
import os

def get_ip_from_mac(target_mac, ip_range="192.168.1.0/24"):
    """
    Finds the IP address associated with a given MAC address on the local network.

    Args:
        target_mac (str): The MAC address to search for (e.g., "00:11:22:33:44:55").
        ip_range (str): The IP address range to scan (e.g., "192.168.1.0/24").

    Returns:
        str: The IP address if found, otherwise None.
    """
    try:
        # Craft an ARP request packet
        # Ether(dst="ff:ff:ff:ff:ff:ff") sends the packet to all devices (broadcast)
        # ARP(pdst=ip_range) targets the specified IP range
        ans, unans = srp(Ether(dst="ff:ff:ff:ff:ff:ff")/ARP(pdst=ip_range), timeout=2, verbose=0)

        for sent, received in ans:
            if received.hwsrc.lower() == target_mac.lower():
                return received.psrc
        return None
    except Exception as e:
        print(f"An error occurred: {e}", file=sys.stderr)
        return None

if __name__ == "__main__":
    macs_to_find = dict()
    mac_list_file_name = "python/mac.txt"
    print(os.getcwd())
    
    try:
        with open(mac_list_file_name, 'r') as file:
            for line in file:
                mac = line.strip()
                macs_to_find[mac] = None
                print(f"Added {mac} to search")
    except FileNotFoundError:
        print("Error: The file '" + mac_list_file_name + "' was not found.")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
    
    # Adjust the IP range to match your local network
    network_range = "10.162.16.1/20" 
    for key in macs_to_find.keys():
        macs_to_find[key] = get_ip_from_mac(key, network_range)

    for i, (key, found_ip) in enumerate(macs_to_find.items()):
        if found_ip:
            print(f"The IP address for MAC {key} is: {found_ip}")
        else:
            print(f"Could not find an IP address for MAC {key} in the range {network_range}.")