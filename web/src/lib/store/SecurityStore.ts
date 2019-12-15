import agent from 'superagent';
import {SecurityModel, ISecurity} from '../model/SecurityModel';
import {indexById, compareByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, observable} from 'mobx';

const query = '{securities {id name scale symbol type version}}';

interface ISecurityResponse {
    body: {data: {securities: ISecurity[]}}
}

const loadingSecurities = 'Loading securities...';

export interface ISecurityStore {
    getSecurity: (id: string | number) => SecurityModel;
    loadSecurities: () => void;
}

export class SecurityStore {
    private loading: boolean = false;
    @observable
    private securitiesById: {[id: string]: SecurityModel} = {};
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    @computed
    get securities(): SecurityModel[] {
        return Object.values(this.securitiesById).sort(compareByName);
    }

    getSecurity(id: string | number): SecurityModel {
        return this.securitiesById['' + id] || {} as SecurityModel;
    }

    loadSecurities(): void {
        if (!this.loading && Object.keys(this.securitiesById).length === 0) {
            this.messageStore.addProgressMessage(loadingSecurities);
            this._loadSecurities();
        }
    }

    private _loadSecurities = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: ISecurityResponse = yield agent.post('/finances/api/v1/graphql').send({query: query});
            this.securitiesById = indexById(data.securities.map(security => new SecurityModel(security)));
        } catch (err) {
            console.error('error gettting securities', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingSecurities);
        }
    });
}