
export async function post<T>(url: string, body: BodyInit): Promise<T> {
    const response = await fetch(url, {method: 'POST', body});
    if (!response.ok) throw new Error(`Fetch failed: ${response.status} - ${response.statusText}`);
    const json = await response.json() as T;
    return json;
}

interface IGraphqlLocation {
    line: number;
    column: number;
}

interface IGraphqlError {
    message: string;
    locations: IGraphqlLocation[];
}

export interface IGraphqlResponse<T> {
    data: T;
    errors?: IGraphqlError[];
}

export async function graphql<T>(url: string, query: string, variables: Record<string, unknown> = {}): Promise<IGraphqlResponse<T>> {
    const body = JSON.stringify({query, variables});
    return post(url, body);
}