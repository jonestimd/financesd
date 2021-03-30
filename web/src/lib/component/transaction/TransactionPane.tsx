import React, {useContext} from 'react';
import {observer} from 'mobx-react-lite';
import TransactionModel, {ITransactionDetail} from '../../model/TransactionModel';
import {Checkbox, TextField} from '@material-ui/core';
import {RootStoreContext} from '../../store/RootStore';
import Autocomplete from '@material-ui/lab/Autocomplete';

interface IProps {
    showSecurity?: boolean;
    transaction?: TransactionModel;
}

const DetailPane = observer<{showSecurity?: boolean, detail: ITransactionDetail}>(({showSecurity, detail}) => {
    const {categoryStore, groupStore} = useContext(RootStoreContext);
    return (
        <div className='detail'>
            <Autocomplete options={categoryStore.categories} getOptionLabel={(c) => c.displayName}
                value={categoryStore.getCategory(detail.transactionCategoryId) ?? null}
                filterOptions={(options, state) => options.filter((c) => c.displayName.toLowerCase().includes(state.inputValue))}
                renderInput={(params) => <TextField {...params} label='Category' variant='outlined' margin='dense' />} />
            <Autocomplete options={groupStore.groups} getOptionLabel={(g) => g.name}
                value={groupStore.getGroup(detail.transactionGroupId) ?? null}
                filterOptions={(options, state) => options.filter((g) => g.name.toLowerCase().includes(state.inputValue))}
                renderInput={(params) => <TextField {...params} label='Group' variant='outlined' margin='dense' />} />
            <TextField label='Memo' variant='outlined' margin='dense' value={detail.memo ?? ''} />
            {showSecurity && <TextField className='number' label='Shares' variant='outlined' margin='dense' value={detail.assetQuantity ?? 0} />}
            <TextField className='number' label='Amount' variant='outlined' margin='dense' value={detail.amount ?? 0} />
        </div>
    );
});

const TransactionPane = observer<IProps>(({showSecurity, transaction}) => {
    const {payeeStore, securityStore} = useContext(RootStoreContext);
    return (
        <div id='transaction-pane'>
            <div className='header'>
                <TextField label='Date' variant='outlined' margin='dense' value={transaction?.date ?? ''} />
                <TextField label='Ref' variant='outlined' margin='dense' value={transaction?.referenceNumber ?? ''} />
                <TextField label='Memo' variant='outlined' margin='dense' value={transaction?.memo ?? ''} className='memo' />
                <Autocomplete options={payeeStore.payees} getOptionLabel={(p) => p.name}
                    value={payeeStore.getPayee(transaction?.payeeId) ?? null}
                    filterOptions={(options, state) => options.filter((p) => p.name.toLocaleLowerCase().includes(state.inputValue))}
                    renderInput={(params) => <TextField {...params} label='Payee' variant='outlined' margin='dense' />} />
                {showSecurity &&
                    <Autocomplete options={securityStore.securities} getOptionLabel={(s) => s.displayName}
                        value={securityStore.getSecurity(transaction?.securityId) ?? null}
                        filterOptions={(options, state) => options.filter((s) => s.displayName.toLocaleLowerCase().includes(state.inputValue))}
                        renderInput={(params) => <TextField {...params} label='Security' variant='outlined' margin='dense' />} />}
                <Checkbox checked={transaction?.cleared ?? false} />
            </div>
            <div className='details'>
                {transaction?.details.map((detail, index) =>
                    <DetailPane key={detail.id ?? -index} showSecurity={showSecurity} detail={detail} />
                )}
            </div>
        </div>
    );
});

export default TransactionPane;
