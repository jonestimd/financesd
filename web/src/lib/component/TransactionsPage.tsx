import React from 'react';
import {Link} from 'react-router-dom';
import TopAppBar from './TopAppBar';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../store/RootStore';
import {translate} from '../i18n/localize';
import Table, {IColumn} from './Table';
import {TransactionModel} from '../model/TransactionModel';
import * as formats from '../formats'

interface IProps {
    match: {params: {[name: string]: string}}
}

const TransactionsPage: React.FC<IProps> = observer(({match: {params: {accountId}}}) => {
    const menuItems = [
        <Link to='/finances/' className='menu-item'>{translate('menu.accounts')}</Link>
    ]
    const {accountStore, payeeStore, transactionStore} = React.useContext(RootStoreContext);
    React.useEffect(() => {
        accountStore.loadAccounts();
        payeeStore.loadPayees();
        transactionStore.loadTransactions(accountId);
    }, []);
    const account = accountStore.getAccount(accountId);
    const columns: IColumn<TransactionModel>[] = [
        {name: translate('transaction.date'), getter: tx => tx.date, className: 'date'},
        {name: translate('transaction.referenceNumber'), getter: tx => tx.referenceNumber},
        {name: translate('transaction.payee'), getter: tx => payeeStore.getPayee(tx.payeeId).name},
        {name: translate('transaction.memo'), getter: tx => tx.memo},
        {name: translate('transaction.cleared'), getter: tx => tx.cleared ? <span>&#x1F5F8;</span> : null, className: 'boolean'},
        {name: translate('transaction.subtotal'), getter: tx => formats.currency.format(tx.subtotal), className: 'number'},
        {name: translate('transaction.balance'), getter: tx => formats.currency.format(tx.balance), className: 'number'},
    ];
    return (
        <>
            <TopAppBar title={account ? account.displayName : ''} menuItems={menuItems} />
            <Table columns={columns} data={transactionStore.getTransactions(accountId)} />
        </>
    );
});

export default TransactionsPage;