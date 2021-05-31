import React from 'react';
import TransactionModel from 'lib/model/TransactionModel';
import {RootStoreContext} from '../../store/RootStore';
import {observer} from 'mobx-react-lite';
import {Icon} from '@material-ui/core';

const Payee: React.FC<{transaction: TransactionModel}> = observer(({transaction: {payeeId}}) => {
    if (payeeId) {
        const {payeeStore} = React.useContext(RootStoreContext);
        return <span className='chip' data-type='payee'><Icon>person</Icon>{payeeStore.getPayee(payeeId)?.name}</span>;
    }
    return null;
});

export default Payee;
