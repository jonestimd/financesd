import {RootStore} from './RootStore';
import * as entityUtils from '../model/entityUtils';
import * as agent from '../agent';
import {newGroup, newGroupModel} from 'test/groupFactory';
import {GroupModel} from '../model/GroupModel';
import {loadingGroups, query} from './GroupStore';

describe('GroupStore', () => {
    const {groupStore, messageStore, alertStore} = new RootStore();

    beforeEach(() => {
        groupStore['groupsById'].clear();
    });
    describe('get groups', () => {
        it('sorts by name', () => {
            const group = newGroupModel();
            groupStore['groupsById'].set(group.id, group);
            jest.spyOn(entityUtils, 'sortValuesByName');

            const groups = groupStore.groups;

            expect(groups).toEqual([group]);
            expect(entityUtils.sortValuesByName).toBeCalledWith(groupStore['groupsById']);
        });
    });
    describe('getGroup', () => {
        it('returns group for ID', () => {
            const group = newGroupModel();
            groupStore['groupsById'].set(group.id, group);

            expect(groupStore.getGroup(group.id)).toBe(group);
        });
        it('returns undefined for unknown ID', () => {
            expect(groupStore.getGroup(-99)).toBeUndefined();
        });
    });
    describe('loadGroups', () => {
        beforeEach(() => {
            groupStore['loading'] = false;
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
        });
        it('loads groups if groupsById is empty', async () => {
            const group = newGroup();
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {groups: [group]}});

            await groupStore.loadGroups();

            expect(groupStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingGroups);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingGroups);
            expect(agent.graphql).toBeCalledWith(query, undefined);
            expect(groupStore.groups).toStrictEqual([new GroupModel(group)]);
        });
        it('does nothing is already loading', async () => {
            groupStore['loading'] = true;
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await groupStore.loadGroups();

            expect(groupStore['loading']).toBe(true);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('does nothing is already loaded', async () => {
            const category = newGroupModel();
            groupStore['groupsById'].set(category.id, category);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await groupStore.loadGroups();

            expect(groupStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });
            jest.spyOn(alertStore, 'addAlert').mockReturnValue();

            await groupStore.loadGroups();

            expect(groupStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingGroups);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingGroups);
            expect(console.error).toBeCalledWith('error from Loading groups', error);
            expect(alertStore.addAlert).toBeCalledWith('error', 'Error loading groups');
        });
    });
});
