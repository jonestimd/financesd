import * as agent from '../agent';
import {SecurityModel, ISecurity} from '../model/SecurityModel';
import {sortValuesByName, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

const query = '{securities {id name scale symbol type version}}';

type SecurityResponse = agent.IGraphqlResponse<{securities: ISecurity[]}>;

const loadingSecurities = 'Loading securities...';

export interface ISecurityStore {
    getSecurity: (id: string | number) => SecurityModel;
    loadSecurities: () => void;
}

export default class SecurityStore {
    private loading = false;
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
            void this._loadSecurities();
        }
    }

    private _loadSecurities = flow(function* (this: SecurityStore): LoadResult<SecurityResponse> {
        this.loading = true;
        try {
            const {data} = yield agent.graphql('/finances/api/v1/graphql', query);
            addToMap(this.securitiesById, data.securities.map((security) => new SecurityModel(security)));
        } catch (err) {
            console.error('error gettting securities', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingSecurities);
        }
    });
}
