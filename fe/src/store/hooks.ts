import { AnyAction, Dispatch, ThunkDispatch } from '@reduxjs/toolkit';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

import { AppDispatch, RootState } from './store';

export const useAppDispatch = (): ThunkDispatch<RootState, null, AnyAction> & Dispatch<AnyAction> =>
	useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
