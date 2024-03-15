## Infer Pipeline 
Infer pipeline, using pipeline to encapsulate inference steps ,such as build samples、recall 、pre-ranking、ranking and rerank.  

Usage:  

   Pipeline(steps=[  
                   ('recall_sample',recallSample),('dssm_recall',dssmRecall),                      // recall  
                   ('pre_ranking_sample','preRankingSample'),('pre_ranking','dssmRank'),           // pre-ranking  
	                ('ranking_sample','rankingSample'),('ranking','deepfm'),                        // ranking  
	                ('re_rank','reRank')                                                            // rerank  
                  ]  
           )  