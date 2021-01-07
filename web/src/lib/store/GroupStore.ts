import * as agent from '../agent';
import {GroupModel, IGroup} from '../model/GroupModel';
import {addToMap, sortValuesByName} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

const query = `{
    groups {id name description version}
}`;

type GroupsResponse = agent.IGraphqlResponse<{groups: IGroup[]}>;

const loadingGroups = 'Loading groups...';

export default class GroupStore {
    private loading = false;
    private groupsById = new ObservableMap<string, GroupModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get groups(): GroupModel[] {
        return sortValuesByName(this.groupsById);
    }

    getGroup(id: string | number): GroupModel {
        return this.groupsById.get('' + id) || {} as GroupModel;
    }

    loadGroups(): void {
        if (!this.loading && this.groupsById.size === 0) {
            this.messageStore.addProgressMessage(loadingGroups);
            void this._loadGroups();
        }
    }

    private _loadGroups = flow(function* (this: GroupStore): LoadResult<GroupsResponse> {
        this.loading = true;
        try {
            const {data} = yield agent.graphql('/finances/api/v1/graphql', query);
            addToMap(this.groupsById, data.groups.map((group) => new GroupModel(group)));
        } catch (err) {
            console.error('error gettting groups', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingGroups);
        }
    });
}
