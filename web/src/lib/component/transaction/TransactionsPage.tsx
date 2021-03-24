import React from 'react';
import {useHistory} from 'react-router';
import TopAppBar from '../TopAppBar';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import {Icon, TextField} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import TransactionList from './TransactionList';
import TransactionTable from './TransactionTable';

interface IProps {
    match: {params: {[name: string]: string}};
}

type ViewMode = 'list' | 'table';

const TransactionsPage: React.FC<IProps> = observer(({match: {params: {accountId: accountParam}}}) => {
    const history = useHistory();
    const accountId = accountParam.length && parseInt(accountParam);
    const {accountStore, transactionStore} = React.useContext(RootStoreContext);
    React.useEffect(() => void transactionStore.loadTransactions(accountId), [transactionStore, accountId]);
    const account = accountStore.getAccount(accountId);
    const [mode, setMode] = React.useState<ViewMode>('list');
    return<>
        <TopAppBar currentPage='transactions'>
            <Autocomplete options={accountStore.accounts} loading={accountStore.accounts.length === 0}
                openOnFocus={true} size='small'
                value={account ?? null} disableClearable={!!account}
                filterSelectedOptions={true}
                groupBy={(account) => account.company?.name ?? ''}
                getOptionLabel={(account) => account.displayName}
                renderInput={(params) => <TextField {...params} variant="outlined" />}
                renderOption={(account) => account.name}
                ListboxProps={{id: 'accounts-menu'}}
                onChange={(_event, value) => value && history.push(`/finances/account/${value.id}`)}/>
            {/* TODO filter input */}
            {/* TODO cleared balance */}
            <ToggleButtonGroup value={mode} exclusive size='small'
                onChange={(_event: React.MouseEvent, value: ViewMode) => setMode(value ?? mode)}>
                <ToggleButton value='list'><Icon>list</Icon></ToggleButton>
                <ToggleButton value='table'><Icon>apps</Icon></ToggleButton>
            </ToggleButtonGroup>
        </TopAppBar>
        {mode === 'list' ? <TransactionList accountId={accountId} /> : <TransactionTable accountId={accountId} />}
    </>;
});

export default TransactionsPage;
