export const fetchIpInfo = (ip: string) => {
    return fetch(`http://211.149.239.251:7777/api/ip/location?ip=${ip}`)
        .then((res) => res.json())
};