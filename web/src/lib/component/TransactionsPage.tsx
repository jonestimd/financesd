import React from 'react';
import {Link} from 'react-router-dom';
import TopAppBar from './TopAppBar';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../store/RootStore';
import {translate} from '../i18n/localize';
import {IColumn} from './Table';
import HeaderDetailTable from './HeaderDetailTable';
import {TransactionModel, ITransactionDetail} from '../model/TransactionModel';
import * as formats from '../formats'

interface IProps {
    match: {params: {[name: string]: string}}
}

const securityAccountTypes = ['BROKERAGE', '_401K'];

const TransactionsPage: React.FC<IProps> = observer(({match: {params: {accountId}}}) => {
    const menuItems = [
        <Link to='/finances/' className='menu-item'>{translate('menu.accounts')}</Link>
    ]
    const {accountStore, categoryStore, payeeStore, transactionStore} = React.useContext(RootStoreContext);
    React.useEffect(() => {
        accountStore.loadAccounts();
        categoryStore.loadCategories();
        payeeStore.loadPayees();
        transactionStore.loadTransactions(accountId);
    }, []);
    const account = accountStore.getAccount(accountId);
    const columns: IColumn<TransactionModel>[] = [
        {key: 'transaction.date', render: tx => tx.date, className: 'date'},
        {key: 'transaction.referenceNumber', render: tx => tx.referenceNumber},
        {key: 'transaction.payee', render: tx => payeeStore.getPayee(tx.payeeId).name},
        {key: 'transaction.memo', render: tx => tx.memo},
        {key: 'transaction.security', render: () => 'tx.securityId', className: 'security'},
        {key: 'transaction.subtotal', render: tx => formats.currency.format(tx.subtotal), className: 'number'},
        {key: 'transaction.cleared', render: tx => tx.cleared ? <span>&#x1F5F8;</span> : null, className: 'boolean'},
        {key: 'transaction.balance', render: tx => formats.currency.format(tx.balance), className: 'number'},
    ];
    const renderCategory = (detail: ITransactionDetail) => {
        if (detail.relatedDetail) return <span className='transfer'>{accountStore.getAccount(detail.relatedDetail.transaction.accountId).name}</span>;
        return <span>{categoryStore.getCategory(detail.transactionCategoryId).displayName}</span>;
    };
    const subcolumns: IColumn<ITransactionDetail>[] = [
        {key: 'detail.group', colspan: 2, render: detail => detail.transactionGroupId},
        {key: 'detail.category', render: renderCategory, className: 'category'},
        {key: 'detail.memo', render: detail => detail.memo},
        {key: 'detail.shares', render: () => 'detail.shares', className: 'security'},
        {key: 'detail.amount', render: detail => formats.currency.format(detail.amount), className: 'number'},
        {key: 'dummy1', header: () => '', render: () => ''},
        {key: 'dummy2', header: () => '', render: () => ''},
    ];
    return (
        <>
            <TopAppBar title={account ? account.displayName : ''} menuItems={menuItems} />
            <HeaderDetailTable
                className={securityAccountTypes.includes(account.type) ? 'security-transactions' : 'transactions'}
                columns={columns}
                subColumns={subcolumns}
                data={transactionStore.getTransactions(accountId)}
                subrows={tx => tx.details} />
        </>
    );
});

export default TransactionsPage;