import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import Memo from './Memo';
import {Checkbox, Typography} from '@material-ui/core';
import ListViewport from '../scroll/ListViewport';
import {useSelection} from '../scroll/selectionHooks';
import Transaction from './Transaction';

interface IProps {
    accountId?: number;
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

const rowSelector = 'div.transaction';
const prototypeSelector = '.prototype';

const TransactionList: React.FC<IProps> = observer(({accountId}) => {
    const {accountStore, transactionStore} = React.useContext(RootStoreContext);
    const tableModel = transactionStore.getTransactionsModel(accountId);
    const showSecurity = accountStore.getAccount(accountId)?.isSecurity;
    const selection = useSelection<HTMLDivElement>({rows: tableModel.transactions.length, rowSelector, prototypeSelector});
    return (
        <ListViewport items={tableModel.transactions} selection={selection}
            renderItem={(tx, _index, selected) => <Transaction tx={tx} selected={selected} showSecurity={showSecurity} />} >
            <TransactionPrototype />
        </ListViewport>
    );
});

export default TransactionList;
