import {GroupModel, IGroup} from '../model/GroupModel';
import {addToMap, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, makeObservable, ObservableMap} from 'mobx';
import Loader from './Loader';
import AlertStore from './AlertStore';

export const query = `{
    groups {id name description version transactionCount}
}`;

export const loadingGroups = 'Loading groups';

export default class GroupStore {
    private loading = false;
    private groupsById = new ObservableMap<number, GroupModel>();
    private loader: Loader;

    constructor(messageStore: IMessageStore, alertStore: AlertStore) {
        makeObservable(this);
        this.loader = new Loader(messageStore, alertStore);
    }

    @computed
    get groups(): GroupModel[] {
        return sortValuesByName(this.groupsById);
    }

    getGroup(id?: number): GroupModel | undefined {
        return typeof id === 'number' ? this.groupsById.get(id) : undefined;
    }

    loadGroups(): Promise<boolean> | undefined {
        if (!this.loading && this.groupsById.size === 0) {
            this.loading = true;
            return this.loader.load<{groups: IGroup[]}>(loadingGroups, {query,
                updater: ({groups}) => addToMap(this.groupsById, groups.map((group) => new GroupModel(group))),
                completer: () => this.loading = false,
            });
        }
    }
}
