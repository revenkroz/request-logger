export interface Log {
    method: string
    url: {
        scheme: string
        host: string
        path: string
        raw_query: string
    }
    full_url: string
    status: string
    status_code: number
    elapsed: number
    req: string
    res: string
    done_at: string

    // frontend
    selected: boolean|undefined
}
