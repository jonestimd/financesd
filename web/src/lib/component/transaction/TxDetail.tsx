import React from 'react';
import {ITransactionDetail} from '../../model/TransactionModel';
import {Currency, Shares} from '../../formats';
import Category from './Category';
import Group from './Group';
import {DetailField} from './Transaction';
import {observer} from 'mobx-react-lite';
import {TextField, TextFieldProps} from '@material-ui/core';
import {translate} from '../../i18n/localize';
import classNames from 'classnames';
import CategoryInput from './CategoryInput';

interface IProps {
    detail: ITransactionDetail;
    editField: DetailField | false;
}

const TxDetail = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, editField, ...inputProps}, ref) => {
    return (
        <div className={classNames('detail chip', {editing: editField})} id={`${detail.id}`}>
            {editField === 'amount'
                ? <TextField ref={ref} {...inputProps} label={translate('Amount')} type='number' value={detail.amount ?? ''} required />
                : <span><Currency>{detail.amount}</Currency></span>}
            {editField === 'category' ? <CategoryInput detail={detail} ref={ref} {...inputProps} /> : <Category detail={detail} />}
            <Group id={detail.transactionGroupId} />
            {detail.assetQuantity ? <span className='shares'><Shares>{detail.assetQuantity}</Shares></span> : null}
            {detail.memo ? <span className='memo'>{detail.memo}</span> : null}
        </div>
    );
}, {forwardRef: true});

export default TxDetail;
