import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import TransactionModel from 'src/lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';
import IconInput from '../IconInput';

interface IProps {
    transaction: TransactionModel;
}

const PayeeInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({transaction, ...inputProps}, ref) => {
    const {payeeStore} = React.useContext(RootStoreContext);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={payeeStore.payees} getOptionLabel={(p) => p.name}
            value={payeeStore.getPayee(transaction.payeeId) ?? null}
            onChange={(_event, value) => transaction.payeeId = value?.id ?? undefined}
            filterOptions={(options, state) => options.filter((p) => p.name.toLocaleLowerCase().includes(state.inputValue.toLowerCase()))}
            renderInput={(params) => <IconInput {...params} {...inputProps} icon='person' />} />
    );

}, {forwardRef: true});

export default PayeeInput;
