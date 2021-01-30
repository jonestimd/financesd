import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table, {IColumn} from './table/Table';
import {PayeeModel} from '../model/PayeeModel';
import {translate} from '../i18n/localize';
import PageMenu from './PageMenu';

const columns: IColumn<PayeeModel>[] = [
    {key: 'payee.name', render: (payee) => payee.name},
    {key: 'payee.transactionCount', render: (payee) => payee.transactionCount, className: 'number'},
];

const PayeesPage: React.FC = observer(() => {
    const {payeeStore} = React.useContext(RootStoreContext);
    const payees = payeeStore.payees;
    return (
        <div className='payee-list'>
            <TopAppBar title={translate('menu.payees')} menuItems={<PageMenu currentPage='menu.payees' />} />
            <Table columns={columns} data={payees} />
        </div>
    );
});

export default PayeesPage;
