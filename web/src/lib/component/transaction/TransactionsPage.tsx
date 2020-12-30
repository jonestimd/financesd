import React from 'react';
import TopAppBar from '../TopAppBar';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import PageMenu from '../PageMenu';
import {Icon} from '@material-ui/core';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import TransactionList from './TransactionList';
import TransactionTable from './TransactionTable';

interface IProps {
    match: {params: {[name: string]: string}};
}

type ViewMode = 'list' | 'table';

const TransactionsPage: React.FC<IProps> = observer(({match: {params: {accountId}}}) => {
    const {accountStore, transactionStore} = React.useContext(RootStoreContext);
    React.useEffect(() => transactionStore.loadTransactions(accountId), [transactionStore, accountId]);
    const account = accountStore.getAccount(accountId);
    const [mode, setMode] = React.useState<ViewMode>('list');
    return (
        <>
            <TopAppBar title={account ? account.displayName : ''} menuItems={<PageMenu />}>
                <ToggleButtonGroup value={mode} exclusive size='small'
                    onChange={(_event: React.MouseEvent, value: ViewMode) => setMode(value ?? mode)}>
                    <ToggleButton value='list'><Icon>list</Icon></ToggleButton>
                    <ToggleButton value='table'><Icon>apps</Icon></ToggleButton>
                </ToggleButtonGroup>
            </TopAppBar>
            {mode === 'list' ? <TransactionList accountId={accountId} /> : <TransactionTable accountId={accountId} />}
        </>
    );
});

export default TransactionsPage;
