import * as agent from './agent';

const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

type ResponseOverrides = Partial<Omit<Response, 'json'>> & {json?: Record<string, unknown>};

function mockResponse({json, ...overrides}: ResponseOverrides): Response {
    return {
        ok: true,
        json: jest.fn().mockResolvedValue(json),
        ...overrides,
    } as unknown as Response;
}

const url = 'some url';
const responseBody = {x: 'the result'};

describe('agent', () => {
    describe('post', () => {
        const requestBody = 'some data';

        it('returns response body', async () => {
            mockFetch.mockResolvedValue(mockResponse({ok: true, json: responseBody}));

            const result = await agent.post(url, requestBody);

            expect(result).toBe(responseBody);
            expect(fetch).toBeCalledWith(url, {method: 'POST', body: requestBody});
        });
        it('throws error for not ok', async () => {
            const statusText = 'request failed';
            mockFetch.mockResolvedValue(mockResponse({ok: false, status: 404, statusText}));

            await expect(agent.post(url, requestBody)).rejects.toThrow(`Fetch failed: 404 - ${statusText}`);
        });
    });
    describe('graphql', () => {
        it('posts query', async () => {
            const query = 'the graphql query';
            mockFetch.mockResolvedValue(mockResponse({json: responseBody}));

            const result = await agent.graphql(query);

            expect(result).toBe(responseBody);
            expect(fetch).toBeCalledWith('/finances/api/v1/graphql', {method: 'POST', body: JSON.stringify({query, variables: {}})});
        });
        it('posts query and variables', async () => {
            const query = 'the graphql query';
            const variables = {a: 123, b: 'some value'};
            mockFetch.mockResolvedValue(mockResponse({json: responseBody}));

            const result = await agent.graphql(query, variables);

            expect(result).toBe(responseBody);
            expect(fetch).toBeCalledWith('/finances/api/v1/graphql', {method: 'POST', body: JSON.stringify({query, variables})});
        });
    });
});
