import * as agent from '../agent';
import {PayeeModel, IPayee} from '../model/PayeeModel';
import {addToMap, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

export const query = `{
    payees {
        id name version transactionCount
    }
}`;

type PayeesResponse = agent.IGraphqlResponse<{payees: IPayee[]}>;

export const loadingPayees = 'Loading payees...';

export default class PayeeStore {
    private loading = false;
    private payeesById = new ObservableMap<string, PayeeModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get payees(): PayeeModel[] {
        return sortValuesByName(this.payeesById);
    }

    getPayee(id: string | number): PayeeModel {
        return this.payeesById.get('' + id);
    }

    loadPayees(): Promise<void> | undefined {
        if (!this.loading && this.payeesById.size === 0) {
            this.messageStore.addProgressMessage(loadingPayees);
            return this._loadPayees();
        }
    }

    private _loadPayees = flow(function* (this: PayeeStore): LoadResult<PayeesResponse> {
        this.loading = true;
        try {
            const {data} = yield agent.graphql('/finances/api/v1/graphql', query);
            addToMap(this.payeesById, data.payees.map((payee) => new PayeeModel(payee)));
        } catch (err) {
            console.error('error gettting payees', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingPayees);
        }
    });
}
