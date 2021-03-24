import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {SecurityModel} from '../model/SecurityModel';
import {translate} from '../i18n/localize';
import {HideZero, Shares} from '../formats';

const columns: IColumn<SecurityModel>[] = [
    {key: 'security.name', render: (security) => security.name},
    {key: 'security.type', className: 'enum', render: (security) => security.type},
    {key: 'security.symbol', render: (security) => security.symbol},
    {key: 'security.shares', className: 'number', render: ({shares}) => shares ? <Shares>{shares}</Shares> : null},
    {key: 'security.firstAcquired', className: 'date', render: (security) => security.firstAcquired},
    {key: 'security.costBasis', className: 'number', render: ({shares, costBasis}) => shares ? <HideZero>{costBasis}</HideZero> : null},
    {key: 'security.dividends', className: 'number', render: (security) => <HideZero>{security.dividends}</HideZero>},
    {key: 'security.transactionCount', render: (security) => security.transactionCount, className: 'number'},
];

const SecuritiesPage: React.FC = observer(() => {
    const {securityStore} = React.useContext(RootStoreContext);
    const securities = securityStore.securities;
    return <>
        <TopAppBar title={translate('menu.securities')} currentPage='menu.securities' />
        <Table columns={columns} data={securities} />
    </>;
});

export default SecuritiesPage;
