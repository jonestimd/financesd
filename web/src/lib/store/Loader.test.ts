import Loader from './Loader';
import {RootStore} from './RootStore';
import * as agent from '../agent';

describe('Loader', () => {
    const {messageStore} = new RootStore();
    const loader = new Loader(messageStore);

    beforeEach(() => {
        jest.spyOn(messageStore, 'addProgressMessage').mockReturnValue();
        jest.spyOn(messageStore, 'removeProgressMessage').mockReturnValue();
    });
    describe('load', () => {
        const message = 'loading...';
        const query = '{ entities {id name} }';
        const variables = {x: 123};
        const updater = jest.fn();
        const completer = jest.fn();

        beforeEach(() => {
            jest.spyOn(console, 'error').mockImplementation(() => {});
        });
        it('calls updater', async () => {
            const result = {entities: []};
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: result});

            expect(await loader.load(message, {query, variables, updater})).toBe(true);

            expect(updater).toBeCalledWith(result);
            expect(agent.graphql).toBeCalledWith(query, variables);
            expect(messageStore.addProgressMessage).toBeCalledWith(message);
            expect(messageStore.removeProgressMessage).toBeCalledWith(message);
        });
        it('calls comleter', async () => {
            const result = {entities: []};
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: result});

            expect(await loader.load(message, {query, variables, updater, completer})).toBe(true);

            expect(updater).toBeCalledWith(result);
            expect(completer).toBeCalledTimes(1);
            expect(agent.graphql).toBeCalledWith(query, variables);
        });
        it('calls completer for graphql errors', async () => {
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: null, errors: [{message: 'syntax error', locations: []}]});

            expect(await loader.load(message, {query, variables, updater, completer})).toBe(false);

            expect(updater).not.toBeCalledTimes(1);
            expect(completer).toBeCalledTimes(1);
            expect(messageStore.addProgressMessage).toBeCalledWith(message);
            expect(messageStore.removeProgressMessage).toBeCalledWith(message);
        });
        it('calls completer after agent error', async () => {
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error('network error'));

            expect(await loader.load(message, {query, variables, updater, completer})).toBe(false);

            expect(updater).not.toBeCalledTimes(1);
            expect(completer).toBeCalledTimes(1);
            expect(messageStore.addProgressMessage).toBeCalledWith(message);
            expect(messageStore.removeProgressMessage).toBeCalledWith(message);
        });
    });
});
