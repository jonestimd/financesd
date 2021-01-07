import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import * as formats from '../../formats';
import classNames from 'classnames';
import TxDetail from './TxDetail';
import Payee from './Payee';
import Security from './Security';
import Memo from './Memo';
import TransactionModel from 'src/lib/model/TransactionModel';
import {Checkbox, Typography} from '@material-ui/core';
import ListViewport from '../scroll/ListViewport';

interface IProps {
    accountId?: string;
}

const TransactionPrototype: React.FC = () => (
    <Typography className='transaction prototype'>
        <div className='leading'><span className='date'>0000-00-00</span></div>
        <div className='details'><Memo text='Prototype' /></div>
        <div className='trailing'>
            <Checkbox checked disabled />
            <span className='number'>0.00</span>
            <span className='number'>0.00</span>
        </div>
    </Typography>
);

const Transaction: React.FC<{tx: TransactionModel, selected: boolean}> = observer(({tx, selected}) => (
    <Typography className={classNames('transaction', {selected})}>
        <div className='leading'>
            <span className='date'>{tx.date}</span>
            {tx.referenceNumber ? <span className='ref-number'>{tx.referenceNumber}</span> : null}
        </div>
        <div className='details'>
            <Payee transaction={tx} />
            <Security transaction={tx} />
            <Memo text={tx.memo} />
            {tx.details.map((detail) => <TxDetail key={detail.id} detail={detail} />)}
        </div>
        <div className='trailing'>
            <Checkbox disabled checked={tx.cleared} />
            <span className={formats.numberClass(tx.subtotal)}>{formats.currency.format(tx.subtotal)}</span>
            <span className={formats.numberClass(tx.balance)}>{formats.currency.format(tx.balance)}</span>
        </div>
    </Typography>
));

const TransactionList: React.FC<IProps> = observer(({accountId}) => {
    const {transactionStore} = React.useContext(RootStoreContext);
    const tableModel = transactionStore.getTransactionsModel(accountId);
    return (
        <ListViewport items={tableModel.transactions} rowSelector='div.transaction' prototypeSelector='.prototype'
            renderItem={(tx, _index, selected) => <Transaction tx={tx} selected={selected} />} >
            <TransactionPrototype />
        </ListViewport>
    );
});

export default TransactionList;
