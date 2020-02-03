import React from 'react';
import TopAppBar from '../TopAppBar';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import * as formats from '../../formats';
import classNames from 'classnames';
import PageMenu from '../PageMenu';
import TxDetail from './TxDetail';
import Payee from './Payee';
import Security from './Security';
import Memo from './Memo';

interface IProps {
    match: {params: {[name: string]: string}};
}

function numberClass(value: number, classes?: string) {
    return classNames(classes, 'number', {negative: value < 0});
}

const TransactionsPage: React.FC<IProps> = observer(({match: {params: {accountId}}}) => {
    const {accountStore, categoryStore, payeeStore, securityStore, transactionStore} = React.useContext(RootStoreContext);
    React.useEffect(() => {
        accountStore.loadAccounts();
        categoryStore.loadCategories();
        payeeStore.loadPayees();
        securityStore.loadSecurities();
        transactionStore.loadTransactions(accountId);
    }, []);
    const account = accountStore.getAccount(accountId);
    return (
        <>
            <TopAppBar title={account ? account.displayName : ''} menuItems={<PageMenu />} />
            <div className='scroll-container'>
                {transactionStore.getTransactionsModel(accountId).transactions.map(tx => (
                    <div key={tx.id} className='transaction'>
                        <div className='leading'>
                            <span className='date'>{tx.date}</span>
                            {tx.referenceNumber ? <span className='ref-number'>{tx.referenceNumber}</span> : null}
                        </div>
                        <div className='details'>
                            <Payee transaction={tx} />
                            <Security transaction={tx} />
                            <Memo text={tx.memo} />
                            {tx.details.map(detail => <TxDetail key={detail.id} detail={detail} />)}
                        </div>
                        <div className='trailing'>
                            <input type='checkbox' checked={tx.cleared} readOnly />
                            <span className={numberClass(tx && tx.subtotal)}>{formats.currency.format(tx.subtotal)}</span>
                            <span className={numberClass(tx && tx.balance)}>{formats.currency.format(tx.balance)}</span>
                        </div>
                    </div>
                ))}
            </div>
        </>
    );
});

export default TransactionsPage;