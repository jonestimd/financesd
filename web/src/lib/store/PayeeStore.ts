import agent from 'superagent';
import {PayeeModel, IPayee} from '../model/PayeeModel';
import {addToMap, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';

const query = `{
    payees {
        id name version
    }
}`;

interface IPayeesResponse {
    body: {data: {payees: IPayee[]}};
}

const loadingPayees = 'Loading payees...';

export default class PayeeStore {
    private loading: boolean = false;
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
        return this.payeesById.get('' + id) || {} as PayeeModel;
    }

    loadPayees(): void {
        if (!this.loading && this.payeesById.size === 0) {
            this.messageStore.addProgressMessage(loadingPayees);
            this._loadPayees();
        }
    }

    private _loadPayees = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: IPayeesResponse = yield agent.post('/finances/api/v1/graphql').send({query});
            addToMap(this.payeesById, data.payees.map(payee => new PayeeModel(payee)));
        } catch (err) {
            console.error('error gettting payees', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingPayees);
        }
    });
}