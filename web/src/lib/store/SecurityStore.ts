import {SecurityModel, ISecurity} from '../model/SecurityModel';
import {sortValuesByName, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, makeObservable, ObservableMap} from 'mobx';
import Loader from './Loader';
import AlertStore from './AlertStore';

export const query = `{
    securities {
        id name type scale symbol type version transactionCount shares firstAcquired costBasis dividends
    }
}`;

export const loadingSecurities = 'Loading securities';

export interface ISecurityStore {
    getSecurity: (id: string | number) => SecurityModel;
    loadSecurities: () => void;
}

export default class SecurityStore {
    private loading = false;
    private securitiesById = new ObservableMap<string, SecurityModel>();
    private loader: Loader;

    constructor(messageStore: IMessageStore, alertStore: AlertStore) {
        makeObservable(this);
        this.loader = new Loader(messageStore, alertStore);
    }

    @computed
    get securities(): SecurityModel[] {
        return sortValuesByName(this.securitiesById);
    }

    getSecurity(id?: string | number): SecurityModel | undefined {
        return this.securitiesById.get('' + id);
    }

    loadSecurities(): Promise<boolean> | undefined {
        if (!this.loading && this.securitiesById.size === 0) {
            this.loading = true;
            return this.loader.load<{securities: ISecurity[]}>(loadingSecurities, {query,
                updater: ({securities}) => addToMap(this.securitiesById, securities.map((security) => new SecurityModel(security))),
                completer: () => this.loading = false,
            });
        }
    }
}
