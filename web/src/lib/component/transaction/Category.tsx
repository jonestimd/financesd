import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import {ITransactionDetail} from 'src/lib/model/TransactionModel';

interface IProps {
    detail: ITransactionDetail;
}

const Category: React.FC<IProps> = observer(({detail: {relatedDetail, transactionCategoryId}}) => {
    const {accountStore, categoryStore} = React.useContext(RootStoreContext);
    if (relatedDetail) {
        return <span className='transfer'>
            <i className='material-icons md-18'>forward</i>
            {accountStore.getAccount(relatedDetail.transaction.accountId).name}
        </span>;
    }
    if (transactionCategoryId) {
        return <span className='category'>{categoryStore.getCategory(transactionCategoryId).displayName}</span>;
    }
    return null;
});

export default Category;
