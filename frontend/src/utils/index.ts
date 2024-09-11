export function extractStatus(statusString: string): number {
    return parseInt(statusString.split(' ')[0])
}
