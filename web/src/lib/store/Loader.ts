import {IMessageStore} from './MessageStore';
import {graphql, IGraphqlResponse} from '../agent';
import {flow} from 'mobx';

type LoadResult<T> = Generator<Promise<IGraphqlResponse<T>>, boolean, IGraphqlResponse<T>>;

interface IOptions<T> {
    query: string;
    variables?: unknown;
    updater: (result: T) => void;
    completer?: () => void;
}

export default class Loader {
    private readonly messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    load = flow(function*<T>(this: Loader, message: string, {query, variables, updater, completer}: IOptions<T>): LoadResult<T> {
        this.messageStore.addProgressMessage(message);
        try {
            const {data, errors} = yield graphql<T>(query, variables);
            if (errors) throw new Error(`graphql errors: ${JSON.stringify(errors)}`);
            updater(data);
            return true;
        } catch (err) {
            console.error('error from ' + message, err);
            return false;
        } finally {
            this.messageStore.removeProgressMessage(message);
            completer?.();
        }
    });
}
