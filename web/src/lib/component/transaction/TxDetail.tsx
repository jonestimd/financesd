import React from 'react';
import {Currency, Shares} from '../../formats';
import Category from './Category';
import Group from './Group';
import {DetailField, securityDetailFields} from '../../model/TransactionModel';
import {observer} from 'mobx-react-lite';
import {TextField, TextFieldProps} from '@material-ui/core';
import CategoryInput from './CategoryInput';
import classNames from 'classnames';
import GroupInput from './GroupInput';
import IconInput from '../IconInput';
import DetailModel from 'lib/model/DetailModel';

interface IProps {
    detail: DetailModel;
    editField: DetailField | false;
    // TODO use to filter categories
    showSecurity?: boolean;
}

const fieldRenderers: Record<DetailField, (detail: DetailModel) => React.ReactNode> = {
    amount: ({amount}) => <span key='amount'><Currency>{amount}</Currency></span>,
    category: (detail) => <Category key='category' detail={detail} />,
    group: ({transactionGroupId}) => <Group key='group' id={transactionGroupId} />,
    shares: ({assetQuantity}) => assetQuantity ? <span key='shares' className='shares'><Shares>{assetQuantity}</Shares></span> : null,
    memo: ({memo}) => memo ? <span key='memo' className='memo'>{memo}</span> : null,
};

const leading = ({detail, editField}: IProps) => {
    if (editField === securityDetailFields[0]) return null;
    const fields = editField ? securityDetailFields.slice(0, securityDetailFields.indexOf(editField)) : securityDetailFields;
    return (
        <div className={classNames('detail chip', {prefix: editField})} id={`${detail.id}`}>
            {fields.map((field: DetailField) => fieldRenderers[field](detail))}
        </div>
    );
};

const trailing = ({detail, editField}: IProps) => {
    if (!editField || editField === securityDetailFields[securityDetailFields.length-1]) return null;
    const fields = securityDetailFields.slice(securityDetailFields.indexOf(editField)+1);
    return (
        <div className='detail chip suffix' id={`${detail.id}`}>
            {fields.map((field) => fieldRenderers[field](detail))}
            <span>&nbsp;</span>
        </div>
    );
};

const TxDetail = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, editField, showSecurity, ...inputProps}) => {
    const {amountText, assetQuantityText, memo} = detail;
    return <>
        {leading({detail, editField})}
        {editField === 'amount' && <TextField {...inputProps} type='number' value={amountText} required
            onChange={({currentTarget}) => detail.amountText = currentTarget.value}
            InputProps={{autoFocus: true, startAdornment: <span>$</span>}} />}
        {editField === 'category' && <CategoryInput detail={detail} {...inputProps} />}
        {editField === 'group' && <GroupInput detail={detail} {...inputProps} />}
        {editField === 'shares' && <IconInput {...inputProps} type='number' value={assetQuantityText} icon='request_page'
            onChange={({currentTarget}) => detail.assetQuantityText = currentTarget.value} />}
        {editField === 'memo' && <IconInput {...inputProps} value={memo ?? ''} icon='notes'
            onChange={({currentTarget}) => detail.memo = currentTarget.value || undefined} />}
        {trailing({detail, editField})}
    </>;
}, {forwardRef: true});

export default TxDetail;
