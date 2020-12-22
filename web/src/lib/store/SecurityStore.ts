import agent from 'superagent';
import {SecurityModel, ISecurity} from '../model/SecurityModel';
import {sortValuesByName, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';

const query = '{securities {id name scale symbol type version}}';

interface ISecurityResponse {
    body: {data: {securities: ISecurity[]}};
}

const loadingSecurities = 'Loading securities...';

export interface ISecurityStore {
    getSecurity: (id: string | number) => SecurityModel;
    loadSecurities: () => void;
}

export default class SecurityStore {
    private loading: boolean = false;
    private securitiesById = new ObservableMap<string, SecurityModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get securities(): SecurityModel[] {
        return sortValuesByName(this.securitiesById);
    }

    getSecurity(id: string | number): SecurityModel {
        return this.securitiesById.get('' + id) || {} as SecurityModel;
    }

    loadSecurities(): void {
        if (!this.loading && this.securitiesById.size === 0) {
            this.messageStore.addProgressMessage(loadingSecurities);
            this._loadSecurities();
        }
    }

    private _loadSecurities = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: ISecurityResponse = yield agent.post('/finances/api/v1/graphql').send({query});
            addToMap(this.securitiesById, data.securities.map(security => new SecurityModel(security)));
        } catch (err) {
            console.error('error gettting securities', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingSecurities);
        }
    });
}