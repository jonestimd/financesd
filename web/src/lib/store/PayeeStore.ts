import {PayeeModel, IPayee} from '../model/PayeeModel';
import {addToMap, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, makeObservable, ObservableMap} from 'mobx';
import Loader from './Loader';
import AlertStore from './AlertStore';

export const query = `{
    payees {
        id name version transactionCount
    }
}`;

export const loadingPayees = 'Loading payees';

export default class PayeeStore {
    private loading = false;
    private payeesById = new ObservableMap<string, PayeeModel>();
    private loader: Loader;

    constructor(messageStore: IMessageStore, alertStore: AlertStore) {
        makeObservable(this);
        this.loader = new Loader(messageStore, alertStore);
    }

    @computed
    get payees(): PayeeModel[] {
        return sortValuesByName(this.payeesById);
    }

    getPayee(id?: string | number): PayeeModel | undefined {
        return this.payeesById.get('' + id);
    }

    loadPayees(): Promise<boolean> | undefined {
        if (!this.loading && this.payeesById.size === 0) {
            this.loading = true;
            return this.loader.load<{payees: IPayee[]}>(loadingPayees, {query,
                updater: ({payees}) => addToMap(this.payeesById, payees.map((payee) => new PayeeModel(payee))),
                completer: () => this.loading = false,
            });
        }
    }
}
