import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextField, TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import {translate} from '../../i18n/localize';
import TransactionModel from 'src/lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';

interface IProps {
    transaction: TransactionModel;
}

const PayeeInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({transaction, ...inputProps}, ref) => {
    const {payeeStore} = React.useContext(RootStoreContext);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={payeeStore.payees} getOptionLabel={(p) => p.name}
            value={payeeStore.getPayee(transaction.payeeId) ?? null}
            filterOptions={(options, state) => options.filter((p) => p.name.toLocaleLowerCase().includes(state.inputValue))}
            renderInput={(params) => <TextField {...params} {...inputProps} label={translate('Payee')} />} />
    );

}, {forwardRef: true});

export default PayeeInput;
