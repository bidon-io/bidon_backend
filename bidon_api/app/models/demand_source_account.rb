class DemandSourceAccount < Sequel::Model
  many_to_one :demand_source
end
