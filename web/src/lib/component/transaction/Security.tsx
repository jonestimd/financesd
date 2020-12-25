import React from 'react';
import TransactionModel from '../../model/TransactionModel';
import {RootStoreContext} from '../../store/RootStore';
import {observer} from 'mobx-react-lite';

const Security: React.FC<{transaction: TransactionModel}> = observer(({transaction: {securityId}}) => {
    if (securityId) {
        const {securityStore} = React.useContext(RootStoreContext);
        const security = securityStore.getSecurity(securityId);
        // icon options: leaderboard request_page security library_books insert_drive_file insert_chart local_atm -trending_up
        return <span className='security chip'><i className='material-icons md-18'>request_page</i>{security.name}</span>;
    }
    return null;
});

export default Security;
