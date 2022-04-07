# Intra-Encrypt

##What It Does So Far
- Creates a node on the device which listens on specified port
- Creates a TUN interface to capture traffic
- Associates with and accepts traffic from the TUN interface(have to implement the forwarding)

##To-Do
- Add directions to create routes on linux host to direct traffic through TUN interface
  - Not too difficult and can be scripted  
- Figure out the neighbor discovery mechanism of this system
  - Possibly requires some sort of bootstrap (centralized) node to keep a list of known peers
  - Figure out the DHT - Awesome concept, but tricky to implement.  
- Decide on centralized node for listing known neighbors
- Add better comments...I know it needs them
- Finish implementing stream handling (the forwarding portion)
