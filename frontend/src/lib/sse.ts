export function listenLogs(callback: (log: any) => void) {
    let es = new EventSource('/logs')

    es.addEventListener('log', (event) => {
        callback(JSON.parse(event.data))
    })
}
