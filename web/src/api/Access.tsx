
export const fetchAccess = () => {
    return fetch(`http://211.149.239.251:7777/api/access`)
        .then((res) => res.json())
};