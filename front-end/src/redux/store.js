import { configureStore } from '@reduxjs/toolkit';
import workLogReducer from './features/workLogSlice'; // Contoh reducer

const store = configureStore({
  reducer: {
    workLog: workLogReducer, // Gabungkan reducer Anda di sini
  },
});

export default store;
