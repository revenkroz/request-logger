export interface Log {
    method: string
    url: {
        Scheme: string
        Host: string
        Path: string
        RawQuery: string
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
