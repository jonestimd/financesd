import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {GroupModel} from '../model/GroupModel';
import {translate} from '../i18n/localize';

const columns: IColumn<GroupModel>[] = [
    {key: 'group.name', render: (group) => group.name},
    {key: 'group.description', render: (group) => group.description},
    {key: 'group.transactionCount', render: (group) => group.transactionCount, className: 'number'},
];

const SecuritiesPage: React.FC = observer(() => {
    const {groupStore} = React.useContext(RootStoreContext);
    const groups = groupStore.groups;
    return <>
        <TopAppBar title={translate('menu.groups')} currentPage='menu.groups' />
        <Table columns={columns} data={groups} />
    </>;
});

export default SecuritiesPage;
