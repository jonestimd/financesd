import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {PayeeModel} from '../model/PayeeModel';
import {translate} from '../i18n/localize';

const columns: IColumn<PayeeModel>[] = [
    {key: 'payee.name', render: (payee) => payee.name},
    {key: 'payee.transactionCount', render: (payee) => payee.transactionCount, className: 'number'},
];

const PayeesPage: React.FC = observer(() => {
    const {payeeStore} = React.useContext(RootStoreContext);
    const payees = payeeStore.payees;
    return <>
        <TopAppBar title={translate('menu.payees')} currentPage='menu.payees' />
        <Table columns={columns} data={payees} />
    </>;
});

export default PayeesPage;
