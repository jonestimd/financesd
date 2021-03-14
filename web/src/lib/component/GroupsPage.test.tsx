import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newGroupModel} from 'src/test/groupFactory';
import GroupsPage from './GroupsPage';
import Table from './table/Table';
import {IColumn} from './table/Column';
import TopAppBar from './TopAppBar';
import {GroupModel} from '../model/GroupModel';
import {ObservableMap} from 'mobx';

describe('GroupsPage', () => {
    const rootStore = new RootStore();
    const {groupStore} = rootStore;
    const group = newGroupModel();

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    it('displays app bar and table of groups', () => {
        const groups = [group];
        jest.spyOn(groupStore, 'groups', 'get').mockReturnValue(groups);

        const component = shallow(<GroupsPage />);

        expect(component.find(TopAppBar)).toHaveProp('title', 'Groups');
        expect(component.find(Table)).toHaveProp('data', groups);
    });
    describe('groups table', () => {
        beforeAll(() => {
            groupStore['groupsById'] = new ObservableMap([
                [group.id, group],
            ]);
        });

        const columnTests: {input: GroupModel, key: string, value: React.ReactNode, name?: string}[] = [
            {input: group, key: 'name', value: group.name},
            {input: group, key: 'description', value: group.description},
            {input: group, key: 'transactionCount', value: group.transactionCount},
        ];

        columnTests.forEach(({name, input, key, value}) => {
            it(name ?? `displays ${key}`, () => {
                const component = shallow(<GroupsPage />);
                const columns = component.find(Table).prop<IColumn<GroupModel>[]>('columns');

                const column = columns.find((column) => column.key === `group.${key}`);

                expect(column?.render(input)).toEqual(value);
            });
        });
    });
});
