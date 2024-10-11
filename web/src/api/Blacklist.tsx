async function all() {
    const response = await fetch(`http://211.149.239.251:7777/api/ip/blacklist`);
    return response.json();
};

async function del(ip: string) {
    const response = await fetch(`http://211.149.239.251:7777/api/ip/blacklist?ip=${ip}`, { method: 'DELETE' });
    return response.json();
};

async function put(ip: string) {
    const response = await fetch(`http://211.149.239.251:7777/api/ip/blacklist?ip=${ip}`, { method: 'PUT' });
    return response.json();
};

export default { all, del, put }