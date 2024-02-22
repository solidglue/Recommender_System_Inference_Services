# all models extend from  basemodel and implement ModelInferInterface
# INFO:每个模型的样本是不一样的，特征也是不一样的，例如长短序列等。需要根据特征构造不同的样本
# 构造器模式？根据离线特征、实时特征、序列特征等构造样本
# 主要区分lr、fm、deepfm类、din类和sim类的区别
