import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table, {IColumn} from './table/Table';
import {SecurityModel} from '../model/SecurityModel';
import {translate} from '../i18n/localize';
import PageMenu from './PageMenu';

const columns: IColumn<SecurityModel>[] = [
    {key: 'security.name', render: (security) => security.name},
    {key: 'security.type', render: (security) => security.type},
    {key: 'security.symbol', render: (security) => security.symbol},
    {key: 'security.transactionCount', render: (security) => security.transactionCount, className: 'number'},
];

const SecuritiesPage: React.FC = observer(() => {
    const {securityStore} = React.useContext(RootStoreContext);
    const securities = securityStore.securities;
    return (
        <div className='securities-list'>
            <TopAppBar title={translate('menu.securities')} menuItems={<PageMenu currentPage='menu.securities' />} />
            <Table columns={columns} data={securities} />
        </div>
    );
});

export default SecuritiesPage;
