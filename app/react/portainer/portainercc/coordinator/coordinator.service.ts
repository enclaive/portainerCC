import axios, { parseAxiosError } from '@/portainer/services/axios';

// export async function getKeys(type: string) {
//     try {
//       const { data } = await axios.get<KeyEntry[]>(buildUrl(), {
//         params: { type },
//       });
//       return data;
//     } catch (error) {
//       throw parseAxiosError(error as Error);
//     }
//   }

//   function buildUrl(id?: KeyId) {
//     let url = '/portainercc/keys';
  
//     if (id) {
//       url += `/${id}`;
//     }
  
//     return url;
//   }