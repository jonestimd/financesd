import agent from 'superagent';
import {PayeeModel, IPayee} from '../model/PayeeModel';
import {indexById, compareByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, observable} from 'mobx';

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
    @observable
    private payeesById: {[id: string]: PayeeModel} = {};
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    @computed
    get payees(): PayeeModel[] {
        return Object.values(this.payeesById).sort(compareByName);
    }

    getPayee(id: string | number): PayeeModel {
        return this.payeesById['' + id] || {} as PayeeModel;
    }

    loadPayees(): void {
        if (!this.loading && Object.keys(this.payeesById).length === 0) {
            this.messageStore.addProgressMessage(loadingPayees);
            this._loadPayees();
        }
    }

    private _loadPayees = flow(function*() {
        this.loading = true;
        try {
            const {body: {data}}: IPayeesResponse = yield agent.post('/finances/api/v1/graphql').send({query});
            this.payeesById = indexById(data.payees.map(payee => new PayeeModel(payee)));
        } catch (err) {
            console.error('error gettting payees', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingPayees);
        }
    });
}